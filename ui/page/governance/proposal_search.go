package governance

import (
	"errors"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/planetdecred/godcr/ui/page/components"
)

// The proposal bleve index directory name.
const (
	proposalIndexDir = "/proposalindex.bleve"
	nameField        = "name"
	usernameField    = "username"
)

// searchResult models the search result.
// It's fields are proposal indexed fields.
type searchResult struct {
	Name     string
	Username string
}

//
func (pg *ProposalsPage) getIndex() (index bleve.Index, err error) {
	if pg.proposalIndex == nil {
		indexPath := pg.WL.Wallet.Root + proposalIndexDir
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
	search.Fields = []string{nameField, usernameField}
	searchResults, err := proposalIndex.Search(search)
	if err != nil {
		log.Error(err)
		return
	}

	var hits []*searchResult
	for _, v := range searchResults.Hits {
		name, ok := v.Fields[nameField].(string)
		if !ok {
			log.Error("Can't assert proposal " + nameField)
		}
		username, ok := v.Fields["username"].(string)
		if !ok {
			log.Error("Can't assert proposal " + usernameField)
		}
		hits = append(hits, &searchResult{
			Name:     name,
			Username: username,
		})
	}
	if len(hits) > 0 {
		pg.proposalSearchChan <- hits
	}
	return
}
