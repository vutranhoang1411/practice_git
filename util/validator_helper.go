package util
const(
	USD="USD"
	GBP="GBP"
	VND="VND"
)
func CheckCurrency(currency string)bool{
	switch currency{
	case USD,VND,GBP:
		return true;
	}
	return false;
}
func CheckEmail (email string)bool{
	count:=0;
	var end string;
	for i:=0;i<len(email);i++{
		if (email[i]!='@'){
			count++;
		}else{
			end=email[i:]
			break;
		}
	}
	return (count>=5)&&(end=="@gmail.com"||end=="@yahoo.com")
}