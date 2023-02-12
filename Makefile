.PHONY: all clean

all: kitchencalendar_no kitchencalendar_us

kitchencalendar_no: calutils.go main.go utils.go nb_NO.go ttf/nunito/Nunito-Bold.ttf ttf/nunito/Nunito-Regular.ttf
	go build -mod=vendor -tags nb_NO -o $@

kitchencalendar_us: calutils.go main.go utils.go en_US.go ttf/nunito/Nunito-Bold.ttf ttf/nunito/Nunito-Regular.ttf
	go build -mod=vendor -tags en_US -o $@

clean:
	rm -f kitchencalendar kitchencalendar_*
