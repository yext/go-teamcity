package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

// ProjectFeature is an interface representing different types of project features that can be added to a project.
type ProjectFeature interface {
	ID() string
	SetID(value string)
	Type() string
	Properties() *Properties
	ProjectID() string
	SetProjectID(value string)
	Disabled() bool
	SetDisabled(value bool)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type projectFeatureJSON struct {
	Disabled   *bool       `json:"disabled,omitempty" xml:"disabled"`
	Href       string      `json:"href,omitempty" xml:"href"`
	ID         string      `json:"id,omitempty" xml:"id"`
	Inherited  *bool       `json:"inherited,omitempty" xml:"inherited"`
	Name       string      `json:"name,omitempty" xml:"name"`
	Properties *Properties `json:"properties,omitempty"`
	Type       string      `json:"type,omitempty" xml:"type"`
}

// ProjectFeatures is a collection of ProjectFeature
type ProjectFeatures struct {
	Count int32                `json:"count,omitempty" xml:"count"`
	Href  string               `json:"href,omitempty" xml:"href"`
	Items []projectFeatureJSON `json:"projectFeature"`
}

// ProjectFeatureService provides operations for managing project features for a project
type ProjectFeatureService struct {
	ProjectID  string
	httpClient *http.Client
	base       *sling.Sling
	restHelper *restHelper
}

func newProjectFeatureService(projectID string, c *http.Client, base *sling.Sling) *ProjectFeatureService {
	slingName := base.New().Path(fmt.Sprintf("projects/%s/projectFeatures/", projectID))
	return &ProjectFeatureService{
		ProjectID:  projectID,
		httpClient: c,
		base:       slingName,
		restHelper: newRestHelperWithSling(c, slingName),
	}
}

// Create adds a new project feature to project
func (s *ProjectFeatureService) Create(pf ProjectFeature) (ProjectFeature, error) {
	if pf == nil {
		return nil, errors.New("pf can't be nil")
	}

	req, err := s.base.New().Post("").BodyJSON(pf).Request()

	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unknown error when adding project feature, statusCode: %d", resp.StatusCode)
	}

	return s.readProjectFeatureResponse(resp)
}

// GetByID returns a project feature by its id
func (s *ProjectFeatureService) GetByID(id string) (ProjectFeature, error) {
	req, err := s.base.New().Get(id).Request()

	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("404 Not Found - Project feature (id: %s) for projectId (id: %s) was not found", id, s.ProjectID)
	}

	return s.readProjectFeatureResponse(resp)
}

// GetProjectFeatures gets all the project features of a Project
func (s *ProjectFeatureService) GetProjectFeatures() ([]ProjectFeature, error) {
	var features ProjectFeatures
	err := s.restHelper.get("", &features, "project features")
	if err != nil {
		return nil, err
	}

	projectFeatures := make([]ProjectFeature, features.Count)

	for i := range features.Items {
		dt, err := json.Marshal(features.Items[i])
		if err != nil {
			return nil, err
		}

		gpf := GenericProjectFeature{}
		err = gpf.UnmarshalJSON(dt)
		if err != nil {
			return nil, err
		}
		projectFeatures[i] = &gpf
	}

	return projectFeatures, nil
}

// Delete removes a project feature from the project configuration by its id.
func (s *ProjectFeatureService) Delete(id string) error {
	request, _ := s.base.New().Delete(id).Request()
	response, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode == 204 {
		return nil
	}

	if response.StatusCode != 200 && response.StatusCode != 204 {
		respData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error '%d' when deleting project feature: %s", response.StatusCode, string(respData))
	}

	return nil
}

// DeleteAll removes all project features of a project configuration
func (s *ProjectFeatureService) DeleteAll() error {
	features, err := s.GetProjectFeatures()
	if err != nil {
		return err
	}

	for _, feature := range features {
		if err := s.Delete(feature.ID()); err != nil {
			return err
		}
	}

	return nil
}

func (s *ProjectFeatureService) readProjectFeatureResponse(resp *http.Response) (ProjectFeature, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload projectFeatureJSON
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, err
	}

	var out ProjectFeature
	var gpf GenericProjectFeature
	if err := gpf.UnmarshalJSON(bodyBytes); err != nil {
		return nil, err
	}
	out = &gpf
	out.SetProjectID(s.ProjectID)
	return out, nil
}
