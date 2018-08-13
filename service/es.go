package service

/*import "github.com/olivere/elastic"

func GetClient() (*elastic.Client, error) {
	var err error

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL("http://"+Config().ES_ip+":"+Config().ES_port+":"),
		elastic.SetBasicAuth(Config().ES_user, Config().ES_password),
		elastic.SetSniff(false))
	if err != nil {
		Exception(err)
	}

	return client, err
}
*/