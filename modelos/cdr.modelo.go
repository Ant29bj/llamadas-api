package modelos

type CDR struct {
	CallDate      string `json:"call_date"`
	Clid          string `json:"clid"`
	Src           string `json:"src"`
	Dst           string `json:"dst"`
	Dcontext      string `json:"dcontext"`
	Channel       string `json:"channel"`
	DstChannel    string `json:"dst_channel"`
	LastApp       string `json:"last_app"`
	LastData      string `json:"last_data"`
	Duration      int    `json:"duration"`
	Billsec       int    `json:"billsec"`
	Disposition   string `json:"disposition"`
	Amaflags      int    `json:"amaflags"`
	Accountcode   string `json:"accountcode"`
	UniqueID      string `json:"unique_id"`
	UserField     string `json:"user_field"`
	Did           string `json:"did"`
	RecordingFile string `json:"recording_file"`
	Cnum          string `json:"cnum"`
	Cnam          string `json:"cnam"`
	OutboundCnum  string `json:"outbound_cnum"`
	OutboundCnam  string `json:"outbound_cnam"`
	DstCnam       string `json:"dst_cnam"`
	LinkedID      string `json:"linked_id"`
	PeerAccount   string `json:"peer_account"`
	Sequence      int    `json:"sequence"`
}
