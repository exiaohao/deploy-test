FROM alpine:3.8

ADD ./tester /usr/local/tester
ADD ./test-data/ /test-data

ENTRYPOINT ["/usr/local/tester", "istio"]
