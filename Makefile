.PHONY: all clean

all: kitchencalendar_no kitchencalendar_us

kitchencalendar_no:
	cd cmd/kitchencalendar && go build -mod=vendor -tags nb_NO -o ../../$@ && cd ../..

kitchencalendar_us:
	cd cmd/kitchencalendar && go build -mod=vendor -tags en_US -o ../../$@ && cd ../..

clean:
	rm -f kitchencalendar kitchencalendar_*
