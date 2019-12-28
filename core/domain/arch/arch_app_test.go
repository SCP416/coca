package arch

import (
	"encoding/json"
	. "github.com/onsi/gomega"
	"github.com/phodal/coca/core/adapter"
	"github.com/phodal/coca/core/adapter/call"
	"github.com/phodal/coca/core/adapter/identifier"
	"github.com/phodal/coca/core/support"
	"io"
	"reflect"
	"testing"
)

func TestConceptAnalyser_Analysis(t *testing.T) {
	g := NewGomegaWithT(t)

	codePath := "../../../_fixtures/arch/step2-java"
	identifierApp := new(identifier.JavaIdentifierApp)
	identifiers := identifierApp.AnalysisPath(codePath)
	var classes []string = nil
	for _, node := range identifiers {
		classes = append(classes, node.Package+"."+node.ClassName)
	}

	callApp := call.NewJavaCallApp()
	callNodes := callApp.AnalysisPath(codePath, classes, identifiers)

	identifiersMap := adapter.BuildIdentifierMap(identifiers)

	app := NewArchApp()
	results := app.Analysis(callNodes, identifiersMap)

	g.Expect(len(results.RelationList)).To(Equal(16))
	g.Expect(len(results.NodeList)).To(Equal(13))

	g.Expect(results.RelationList["domain.AggregateRootA->domain.AggregateRoot"].From).To(Equal("domain.AggregateRootA"))
	g.Expect(results.RelationList["domain.AggregateRootA->domain.AggregateRoot"].To).To(Equal("domain.AggregateRoot"))

	jsonContent, _ := json.MarshalIndent(results, "", "\t")
	content := support.ReadFile(codePath + "/" + "results.json")

	g.Expect(JSONBytesEqual(jsonContent, content)).To(Equal(true))
}

func JSONEqual(a, b io.Reader) (bool, error) {
	var j, j2 interface{}
	d := json.NewDecoder(a)
	if err := d.Decode(&j); err != nil {
		return false, err
	}
	d = json.NewDecoder(b)
	if err := d.Decode(&j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
