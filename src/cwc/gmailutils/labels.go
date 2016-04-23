package gmailutils

import (
	gmail "google.golang.org/api/gmail/v1"
)

func Labels(srv *gmail.Service, user string) (map[string]string, error) {
	o := make(map[string]string)
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		return nil, err
	}
	for _, l := range r.Labels {
		o[l.Name] = l.Id
	}
	return o, nil
}
