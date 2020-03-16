package db

type Index struct {
	Table        string `gorm:"column:Table"`
	NonUnique    bool   `gorm:"column:Non_unique"`
	KeyName      string `gorm:"column:Key_name"`
	Seq          int
	ColumnName   string
	Collation    string
	Cardinality  int
	Sub          string
	Packed       string
	Null         string
	Type         string
	Comment      string
	IndexComment string
}
