FROM golang:1.16-alpine 

WORKDIR /go/src/gitlab.com/spezi/services/pdf_text_extractor

RUN apk update;
RUN apk add --no-cache g++ \
                       tesseract-ocr-dev~=4.1 \
                       tesseract-ocr-data-por~=4.1 \
                       imagemagick6-dev \
                       ghostscript;

COPY policy.xml /etc/ImageMagick-6/

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["pdf_text_extractor"]
