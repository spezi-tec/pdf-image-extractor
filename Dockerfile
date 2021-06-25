FROM golang:1.16 

WORKDIR /go/src/gitlab.com/spezi/services/pdf_text_extractor

RUN apt-get update;
RUN apt-get install -y libtesseract-dev \
                       tesseract-ocr-por \
                       libmagickwand-dev \
                       imagemagick-6.q16 \
                       ghostscript;

COPY policy.xml /etc/ImageMagick-6/

COPY . .

RUN go get -d -v ./...
RUN GOOS=linux go build -a -o pdf_text_extractor .

CMD ["./pdf_text_extractor"]
