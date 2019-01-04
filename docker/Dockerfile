FROM byrnedo/alpine-curl

ADD ./tester /usr/local/tester
ADD ./test-data/ /test-data

ENTRYPOINT ["/usr/local/tester", "istio"]
