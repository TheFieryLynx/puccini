package csar

import (
	contextpkg "context"
	"fmt"
	"strconv"

	"github.com/tliron/exturl"
)

type ServiceTemplateSelection struct {
	URL           exturl.URL
	RepositoryURL exturl.URL
	Meta          *Meta
	MetaPresent   bool
	RootFallback  bool
}

func NewURL(csarUrl exturl.URL, format string, path string) (exturl.URL, exturl.URL, error) {
	if format == "" {
		format = csarUrl.Format()
	}

	if exturl.IsValidTarballArchiveFormat(format) {
		return exturl.NewTarballURL(path, csarUrl, format), exturl.NewTarballURL("", csarUrl, format), nil
	}

	switch format {
	case "zip", "csar":
		return exturl.NewZipURL(path, csarUrl), exturl.NewZipURL("", csarUrl), nil
	default:
		return nil, nil, fmt.Errorf("unsupported CSAR archive format: %q", format)
	}
}

func GetDefaultServiceTemplateURL(context contextpkg.Context, csarUrl exturl.URL, format string) (exturl.URL, exturl.URL, error) {
	return GetServiceTemplateURL(context, csarUrl, format, "")
}

func GetServiceTemplateURL(context contextpkg.Context, csarUrl exturl.URL, format string, serviceTemplateName string) (exturl.URL, exturl.URL, error) {
	selection, err := ResolveServiceTemplateURL(context, csarUrl, format, serviceTemplateName)
	if err != nil {
		return nil, nil, err
	}
	return selection.URL, selection.RepositoryURL, nil
}

func ResolveServiceTemplateURL(context contextpkg.Context, csarUrl exturl.URL, format string, serviceTemplateName string) (*ServiceTemplateSelection, error) {
	if format == "" {
		format = csarUrl.Format()
	}

	metaPresent := true
	meta, err := ReadMetaFromURL(context, csarUrl, format)
	if err != nil {
		if exturl.IsNotFound(err) {
			metaPresent = false
			if meta, err = NewMetaFor(context, csarUrl, format); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if (serviceTemplateName == "") || (serviceTemplateName == "0") {
		// Default entry point

		// Attempt to use Entry-Definitions in TOSCA.meta
		if meta.EntryDefinitions != "" {
			url, repositoryURL, err := NewURL(csarUrl, format, meta.EntryDefinitions)
			if err != nil {
				return nil, err
			}
			return &ServiceTemplateSelection{
				URL:           url,
				RepositoryURL: repositoryURL,
				Meta:          meta,
				MetaPresent:   metaPresent,
				RootFallback:  !metaPresent,
			}, nil
		}

		// Attempt to find it in root of archive
		if path, err := GetRootPath(context, csarUrl, format); err == nil {
			url, repositoryURL, err := NewURL(csarUrl, format, path)
			if err != nil {
				return nil, err
			}
			return &ServiceTemplateSelection{
				URL:           url,
				RepositoryURL: repositoryURL,
				Meta:          meta,
				MetaPresent:   metaPresent,
				RootFallback:  true,
			}, nil
		} else {
			return nil, err
		}
	} else {
		// Alternative entry points

		// Try as integer
		if serviceTemplateNumber, err := strconv.ParseUint(serviceTemplateName, 10, 64); err == nil {
			if otherDefinitionIndex := int(serviceTemplateNumber) - 1; (otherDefinitionIndex >= 0) && (otherDefinitionIndex < len(meta.OtherDefinitions)) {
				url, repositoryURL, err := NewURL(csarUrl, format, meta.OtherDefinitions[otherDefinitionIndex])
				if err != nil {
					return nil, err
				}
				return &ServiceTemplateSelection{
					URL:           url,
					RepositoryURL: repositoryURL,
					Meta:          meta,
					MetaPresent:   metaPresent,
				}, nil
			}
		}

		// Try as string
		for _, otherDefinition := range meta.OtherDefinitions {
			if otherDefinition == serviceTemplateName {
				url, repositoryURL, err := NewURL(csarUrl, format, serviceTemplateName)
				if err != nil {
					return nil, err
				}
				return &ServiceTemplateSelection{
					URL:           url,
					RepositoryURL: repositoryURL,
					Meta:          meta,
					MetaPresent:   metaPresent,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("CSAR does not have service template %q: %s", serviceTemplateName, csarUrl.String())
}
