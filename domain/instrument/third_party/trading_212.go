package third_party

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/instrument_service"

	"github.com/chromedp/chromedp"
	"github.com/dystopia-systems/alaskalog"
	"github.com/gofrs/uuid"
	"github.com/vectorman1/analysis/analysis-api/common"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

type externalSymbolService interface {
	GetLatest(*context.Context) (*[]*instrument_service.Instrument, error)
}

type ExternalSymbolService struct {
}

func NewTrading212Service() *ExternalSymbolService {
	return &ExternalSymbolService{}
}

func (s *ExternalSymbolService) GetLatest(ctx context.Context) (*[]*instrument_service.Instrument, error) {
	bctx, c1 := chromedp.NewContext(
		ctx,
		chromedp.WithLogf(alaskalog.Logger.Infof),
	)
	defer c1()

	tbctx, c2 := context.WithTimeout(bctx, 30*time.Second)
	defer c2()

	var htmlRes string
	err := chromedp.Run(tbctx,
		chromedp.Navigate(common.Trading212InstrumentsLink),
		chromedp.WaitVisible(common.Trading212ShowAllButtonSelector),
		chromedp.Click(common.Trading212ShowAllButtonSelector),
		chromedp.InnerHTML(common.Trading212AllInstrumentsSelector, &htmlRes))
	if err != nil {
		alaskalog.Logger.Warnf("failed to get 212 webpage: %v", err)
		return nil, err
	}

	return parseHtmlToProtoSyms(htmlRes)
}

func parseHtmlToProtoSyms(htmlRes string) (*[]*instrument_service.Instrument, error) {
	var parsedProtoSyms []*instrument_service.Instrument
	var wg sync.WaitGroup

	rows, err := walkTable(htmlRes)
	if err != nil {
		return nil, err
	}

	// get results from parser worker
	parsedSymsChan := make(chan *instrument_service.Instrument)
	go func(wg *sync.WaitGroup) {
		for sym := range parsedSymsChan {
			parsedProtoSyms = append(parsedProtoSyms, sym)
			wg.Done()
		}
	}(&wg)

	ctx, c1 := context.WithTimeout(context.Background(), time.Second)
	defer c1()

	g, _ := errgroup.WithContext(ctx)

	// spawn goroutine for each row
	for _, row := range rows {
		wg.Add(1)
		trow := row
		g.Go(func() error {
			sym, err := getSymbolData(trow)
			if err != nil {
				return err
			}
			parsedSymsChan <- sym
			return nil
		})
	}

	// check for any errors
	err = g.Wait()
	if err != nil {
		return nil, err
	}

	// wait for the results to be added to the array
	wg.Wait()

	return &parsedProtoSyms, nil
}

// getSymbolData reads a row from the table and parses it into a proto struct
func getSymbolData(row []string) (*instrument_service.Instrument, error) {
	instrumentName := strings.TrimSpace(row[0])
	companyName := strings.TrimSpace(row[1])
	currencyCode := strings.TrimSpace(row[2])
	isin := strings.TrimSpace(row[3])
	minTradedQuantity, _ := strconv.ParseFloat(strings.TrimSpace(row[4]), 32)
	roundedMinQuantity := float32(math.Round(minTradedQuantity*1000) / 1000)
	marketName := strings.TrimSpace(row[5])
	marketHours := strings.TrimSpace(row[6])

	ns, err := uuid.FromString(common.SymbolsNamespace)
	if err != nil {
		return nil, err
	}

	str := fmt.Sprintf("%s,%s,%s", isin, instrumentName, marketName)
	u := uuid.NewV5(ns, str)
	us := u.String()

	return &instrument_service.Instrument{
		Uuid:                 us,
		Isin:                 isin,
		Identifier:           instrumentName,
		Name:                 companyName,
		CurrencyCode:         currencyCode,
		MinimumOrderQuantity: roundedMinQuantity,
		MarketName:           marketName,
		MarketHoursGmt:       marketHours,
	}, nil
}

// walkTable recursively walks the table of instruments received and returns it as a splice of splices
func walkTable(htmlRes string) ([][]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlRes))
	if err != nil {
		return nil, err
	}

	var symbolRows [][]string

	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "id" && strings.Contains(a.Val, "equity-row-") {
					var row []string
					instrumentName := n.FirstChild
					companyName := instrumentName.NextSibling
					currencyCode := companyName.NextSibling
					isin := currencyCode.NextSibling
					minTradedQuantity := isin.NextSibling
					marketName := minTradedQuantity.NextSibling
					marketHours := marketName.NextSibling

					row = append(row,
						[]string{
							instrumentName.FirstChild.Data,
							companyName.FirstChild.Data,
							currencyCode.FirstChild.Data,
							isin.FirstChild.Data,
							minTradedQuantity.FirstChild.Data,
							marketName.FirstChild.Data,
							strings.TrimSpace(marketHours.FirstChild.Data)}...)
					symbolRows = append(symbolRows, row)
				}
			}
		}
	}

	forEachNode(doc, visitNode, nil)

	return symbolRows, nil
}

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
