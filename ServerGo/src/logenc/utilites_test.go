package logenc

import (
	"reflect"
	"testing"
	"time"
)

func TestGenTestULID(t *testing.T) {
	tests := []struct {
		name    string
		notwant Log
	}{
		{
			name:    "TestGenTestLogWithULID",
			notwant: Log{XML_ULID: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var listlog Log
			listlog.GenTestULID(time.Now())
			if got := listlog.XML_ULID; reflect.DeepEqual(got, tt.notwant) {
				t.Errorf("GenTestLogWithULID() = %v, want %v", got, tt.notwant)
			}
		})
	}
}

func Test_Datestr2time(t *testing.T) {
	type args struct {
		in0 string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "t1",
			args: args{in0: "08092021224536920"},
			want: time.Date(2021, 9, 8, 22, 45, 36, 920000000, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := datestr2time(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("datestr2time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeLine(t *testing.T) {
	line := "<loglist><log module_name=\" Monitor\" app_path=\"/usr/local//monitor/\" app_pid=\"4913\" thread_id=\"\" time=\"29052021000147040\" ulid=\"0001GB313BF4HPFYCDY3QTZ6A6\" type=\"3\" message=\"Состояние '[192.168.1.120] Cервер ' изменилось на 'Ошибка'\" ext_message=\"Context:  -- void ::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/src/123.cpp:686\"/></loglist>"

	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestEncodeLine",
			args: args{line: line},
			want: "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZb3Z4aBt2VFVST1RJGRtaS0tkS1pPUwYZFE5ISRRXVFhaVxR3XlZBFE9WWEgUVlRVUk9USRRPVlhIZFZUVVJPVEkZG1pLS2RLUl8GGQ8CCggZG09TSV5aX2RSXwYZGRtPUlZeBhkJAgsOCQsJCgsLCwoPDAsPCxkbTldSXwYZCwsLCnx5CAoIeX0Pc2t9Ynh/Yghqb2ENeg0ZG09CS14GGQgZG1ZeSEhaXF4GGeua64Xquuq564XqtOuG64PrjhscYAoCCRUKDQMVChUKCQtmG3jrjuq764nrjuq7G+uh65pkaXdyFOuh65rrqRvrmeuF64Tri+uMHBvrg+uM64frjuuG64PrgOuF6rrqtxvrhuuLGxzrpeqz64PriuuB64scGRteQ09kVl5ISFpcXgYZeFRVT15DTwEbGxYWG01UUl8bT1ZYSAEBellIT0laWE92VFVST1RJAQFUVXhUVktUVV5VT2hPWk9eeFNaVVxeXxNqbk5SXxIAGxUVFBUVFBUVFBUVFEhJWBRXUllIFE9WWEhkS1dOXFJVFEhJWBR6WUhPSVpYT3ZUVVJPVEkVWEtLAQ0DDRkUBQcUV1RcV1JITwU=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeLine(tt.args.line); got != tt.want {
				t.Errorf("EncodeLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeLine(t *testing.T) {
	type args struct {
		line string
	}
	//line3 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZb2hoG2heSU1SWF4ZG1pLS2RLWk9TBhkUTkhJFFdUWFpXFHdeVkEUT0hIFE9ISGRIXklNUlheGRtaS0tkS1JfBhkKDA8JGRtPU0leWl9kUl8GGVVPS2RYVxkbT1JWXgYZCwoLDAkLCQoLCwkOCgMJDwMZG05XUl8GGQsLCwp8f28JbX51eGxjc2IPemgOeWFjb3MMGRtPQkteBhkLGRtWXkhIWlxeBhkdSk5UTwDrrOuG64vqvOuO64brg+uOG+uE64XrgOq0G+q66rnqu+uL6rnquOuHG+uJG3Vvaxvrheq564nrjuq5644b64TrheuA6rjqvOuO64brhuuF64cb64XquRt1b2sb6rrrjuq764nrjuq764sbCgIJFQoNAxUKFQkOCQEKCQgb64PrjOuH647rhuuD64Drheq66rcBGwobFh1cTwAbCgsVHUpOVE8AGRQFBxRXVFxXUkhPBQ=="
	///line1 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZb3Z4aBt2VFVST1RJGRtaS0tkS1pPUwYZFE5ISRRXVFhaVxR3XlZBFE9WWEgUVlRVUk9USRRPVlhIZFZUVVJPVEkZG1pLS2RLUl8GGQ8CCggZG09TSV5aX2RSXwYZGRtPUlZeBhkJAgsOCQsJCgsLCwoPDAsPCxkbTldSXwYZCwsLCnx5CAoIeX0Pc2t9Ynh/Yghqb2ENeg0ZG09CS14GGQgZG1ZeSEhaXF4GGeua64Xquuq564XqtOuG64PrjhscYAoCCRUKDQMVChUKCQtmG3jrjuq764nrjuq7G+uh65pkaXdyFOuh65rrqRvrmeuF64Tri+uMHBvrg+uM64frjuuG64PrgOuF6rrqtxvrhuuLGxzrpeqz64PriuuB64scGRteQ09kVl5ISFpcXgYZeFRVT15DTwEbGxYWG01UUl8bT1ZYSAEBellIT0laWE92VFVST1RJAQFUVXhUVktUVV5VT2hPWk9eeFNaVVxeXxNqbk5SXxIAGxUVFBUVFBUVFBUVFEhJWBRXUllIFE9WWEhkS1dOXFJVFEhJWBR6WUhPSVpYT3ZUVVJPVEkVWEtLAQ0DDRkUBQcUV1RcV1JITwU="
	line2 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZDG92eGgbb35obxkbWktLZEtaT1MGGRQIFG9+aG8Ub35obxkbWktLZEtSXwYZCQILGRtPU0leWl9kUl8GGQkZG09SVl4GGQkCCw4JCwkKCwsLCg8MCw8LGRtOV1JfBhkLCn11eWlxDHkDcX5xYnZ9fwMJfQtvfX94CxkbT0JLXgYZCRkbVl5ISFpcXgYZ65rrheq66rnrheq064brg+uOGxwCAxUKCQoVCgMKFQoJA3jrjuq764nrjuq7G+uh65pkbn9rFOuk64PrhuuIHBkbXkNPZFZeSEhaXF4GGXhUVU9eQ08BGxsWFhtNVFJfG09WWEgBAXpZSE9JWlhPdlRVUk9USQEBF3xJXlpPXkkbd1RVX1RVGQUHFFdUXFdSSE8F"
	//b64data := line[strings.IndexByte(line, ',')+1:]
	//line2 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZSV5aX1JVS05PGRtaS0tkS1pPUwYZFFNUVl4UTkheSQsUT1RLWkEUSV5aX1JVS05PFEleWl9SVUtOTxkbWktLZEtSXwYZDQ8NAxkbT1NJXlpfZFJfBhkKCAICAw0CCgwNCAICCA0ZG09SVl4GGQoMCw4JCwkKCwsLCwsLCw8LGRtPQkteBhkIGRtWXkhIWlxeBhl+aWl0aQEbG+um644b6rrrheuJ64Tri+uP64vrjuq5G+q864PquuuA64Ub64Tqu+uF6rzrg+q564vrhuuG6rDqvhvriuuL64LquRvrgxvqvOuD6rrrgOuFG+uK64vrguq5G+q664XrheuK6rLrjuuG64PqtBsbeHR1dX54b3J0dQEbCA8ba2l0b3R4dHcBG2AKCWYbekhPXklSQxPrqOuk65sREhvrm+uL64/ri+q76rABGwgMGxsKDBULDhUJCwkKGwsLAQsLAQsLGxsKCxsLCRsJCRsLCxsLeRt9CxsOChsLCxsLCRsMCxt4CRsIfxsDCxsLfxsICxsKCxsLCBsZFAUHFFdUXFdSSE8F"
	//line2 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZSV5aX1JVS05PGRtaS0tkS1pPUwYZFFNUVl4UTkheSQsUT1RLWkEUSV5aX1JVS05PFEleWl9SVUtOTxkbWktLZEtSXwYZDQ8NAxkbT1NJXlpfZFJfBhkKCAICAw0CCgwNCAICCA0ZG09SVl4GGQoMCw4JCwkKCwsLCwsLCw8LGRtPQkteBhkIGRtWXkhIWlxeBhl+aWl0aQEbG+um644b6rrrheuJ64Tri+uP64vrjuq5G+q864PquuuA64Ub64Tqu+uF6rzrg+q564vrhuuG6rDqvhvriuuL64LquRvrgxvqvOuD6rrrgOuFG+uK64vrguq5G+q664XrheuK6rLrjuuG64PqtBsbeHR1dX54b3J0dQEbCA8ba2l0b3R4dHcBG2AKCWYbekhPXklSQxPrqOuk65sREhvrm+uL64/ri+q76rABGwgMGxsKDBULDhUJCwkKGwsLAQsLAQsLGxsKCxsLCRsJCRsLCxsLeRt9CxsOChsLCxsLCRsMCxt4CRsIfxsDCxsLfxsICxsKCxsLCBsZFAUHFFdUXFdSSE8F"
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestDecodeLine",
			args: args{line: line2}, //"01FNBRJ7B8JEJYMFD82F0TFDC0"
			want: "<loglist><log module_name=\" TEST\" app_path=\"/3/TEST/TEST\" app_pid=\"290\" thread_id=\"2\" time=\"29052021000147040\" ulid=\"01FNBRJ7B8JEJYMFD82F0TFDC0\" type=\"2\" message=\"Состояние '98.121.181.128Cервер КС_UDP/Пинг'\" ext_message=\"Context:  -- void ::AbstractMonitor::,Greater London\"></loglist>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeLine(tt.args.line); got != tt.want {
				t.Errorf("DecodeLine() = %v, want %v", got, tt.want)

			}
		})
	}
}
