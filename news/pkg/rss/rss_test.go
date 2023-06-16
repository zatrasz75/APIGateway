package rss

import (
	storage "GoNews/pkg/storage"
	"testing"
)

func TestGetRss(t *testing.T) {
	posts, err := GetRss("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(posts), posts)

	posts, err = GetRss("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(posts), posts)

	posts, err = GetRss("https://cprss.s3.amazonaws.com/golangweekly.com.xml")
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(posts), posts)
}

func TestGoNews(t *testing.T) {
	chP := make(chan []storage.Post)
	chE := make(chan error)

	type args struct {
		configURL string
		chP       chan []storage.Post
		chE       chan error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{configURL: "./config.json", chP: chP, chE: chE},
			wantErr: false,
		},
		{
			name:    "null config",
			args:    args{configURL: "", chP: chP, chE: chE},
			wantErr: true,
		},
		{
			name:    "invalid config",
			args:    args{configURL: "./test_config_invalid.json", chP: chP, chE: chE},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GoNews(tt.args.configURL, tt.args.chP, tt.args.chE); (err != nil) != tt.wantErr {
				t.Errorf("GoNews() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
