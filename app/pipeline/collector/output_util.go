package collector

import (
	"github.com/andeya/pholcus/logs"
)

// namespace returns the main namespace (relative to DB name); optional, does not depend on data content.
func (self *Collector) namespace() string {
	if self.Spider.Namespace == nil {
		if self.Spider.GetSubName() == "" {
			return self.Spider.GetName()
		}
		return self.Spider.GetName() + "__" + self.Spider.GetSubName()
	}
	return self.Spider.Namespace(self.Spider)
}

// subNamespace returns the sub-namespace (relative to table name); optional, may depend on data content.
func (self *Collector) subNamespace(dataCell map[string]interface{}) string {
	if self.Spider.SubNamespace == nil {
		return dataCell["RuleName"].(string)
	}
	defer func() {
		if p := recover(); p != nil {
			logs.Log.Error("subNamespace: %v", p)
		}
	}()
	return self.Spider.SubNamespace(self.Spider, dataCell)
}

// joinNamespaces concatenates main and sub-namespace with double underscore.
func joinNamespaces(namespace, subNamespace string) string {
	if namespace == "" {
		return subNamespace
	} else if subNamespace != "" {
		return namespace + "__" + subNamespace
	}
	return namespace
}
