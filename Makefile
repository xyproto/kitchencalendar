.PHONY: all clean

all: kitchencalendar_no kitchencalendar_us

kitchencalendar_no: main.go nb_NO.go
	go build -mod=vendor -tags nb_NO -o $@

kitchencalendar_us: main.go en_US.go
	go build -mod=vendor -tags en_US -o $@

clean:
	rm kitchencalendar_*
