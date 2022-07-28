package governance

import (
	"errors"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/highlight/highlighter/ansi"

	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

// The proposal bleve index directory name.
const (
	proposalIndexDir = "/proposalindex.bleve"
	nameField        = "name"
	usernameField    = "username"
	tokenField       = "token"
)

// searchResult models the search result.
// It's fields are proposal indexed fields.
type searchResult struct {
	Token    string
	Name     string
	Username string
}

func (pg *ProposalsPage) getIndex() (index bleve.Index, err error) {
	if pg.proposalIndex == nil {
		indexPath := pg.WL.WalletDirectory() + proposalIndexDir
		if indexPath == "" {
			return index, errors.New("Error: Empty config root")
		}
		index, err = bleve.Open(indexPath)
		if err == bleve.ErrorIndexPathDoesNotExist {
			var indexMapping mapping.IndexMapping
			indexMapping, err = pg.buildIndexMapping()
			if err != nil {
				return
			}

			index, err = bleve.New(indexPath, indexMapping)
			if err != nil {
				return
			}
		}

		pg.proposalIndex = index
		return
	} else {
		return pg.proposalIndex, nil
	}
}

func (pg *ProposalsPage) indexProposal(proposals []*components.ProposalItem) error {
	proposalIndex, err := pg.getIndex()
	if err != nil {
		return err
	}

	// index proposals.
	go func() {
		err = pg.createProposalIndex(proposalIndex, proposals)
		if err != nil {
			log.Error(err)
		}
		return
	}()

	return nil
}

func (pg *ProposalsPage) createProposalIndex(propIndex bleve.Index, proposals []*components.ProposalItem) error {
	for _, prop := range proposals {
		err := propIndex.Index(prop.Proposal.Name, prop.Proposal) // index proposal using proposal name as ID.
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (pg *ProposalsPage) buildIndexMapping() (mapping.IndexMapping, error) {
	// a generic reusable mapping for english word text
	engTextFieldMapping := bleve.NewTextFieldMapping()
	engTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	proposalMapping := bleve.NewDocumentMapping()
	proposalMapping.AddFieldMappingsAt(nameField, engTextFieldMapping)
	proposalMapping.AddFieldMappingsAt(usernameField, keywordFieldMapping)
	proposalMapping.AddFieldMappingsAt(tokenField, keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("proposal", proposalMapping)

	return indexMapping, nil
}

func (pg *ProposalsPage) searchProposal(searchTerm string) {
	proposalIndex, err := pg.getIndex()
	if err != nil {
		log.Info(err)
		return
	}
	query := bleve.NewMatchQuery(searchTerm)
	search := bleve.NewSearchRequest(query)
	search.Highlight = bleve.NewHighlightWithStyle(ansi.Name)
	search.Fields = []string{nameField, usernameField, tokenField}
	searchResults, err := proposalIndex.Search(search)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Search Result: %v", searchResults)

	var hits []*searchResult
	for _, v := range searchResults.Hits {
		name, ok := v.Fields[nameField].(string)
		if !ok {
			log.Error("Can't assert proposal " + nameField)
		}
		username, ok := v.Fields[usernameField].(string)
		if !ok {
			log.Error("Can't assert proposal " + usernameField)
		}
		token, ok := v.Fields[tokenField].(string)
		if !ok {
			log.Error("Can't assert proposal " + tokenField)
		}
		hits = append(hits, &searchResult{
			Token:    token,
			Name:     name,
			Username: username,
		})
	}

	pg.proposalSearchChan <- hits
	return
}

func (pg *ProposalsPage) layoutSearchResult(gtx C, l *load.Load, result *searchResult) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.layoutAuthor(gtx, l, result)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.layoutTitle(gtx, l, result)
			}),
		)
	})
}

func (pg *ProposalsPage) layoutTitle(gtx C, l *load.Load, result *searchResult) D {
	lbl := l.Theme.H6(result.Name)
	lbl.Font.Weight = text.SemiBold
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func (pg *ProposalsPage) layoutAuthor(gtx C, l *load.Load, result *searchResult) D {
	grayCol := l.Theme.Color.GrayText2
	nameLabel := l.Theme.Body2(result.Username)
	nameLabel.Color = grayCol

	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(nameLabel.Layout),
			)
		}),
	)
}
