FROM scratch
CMD ["/zenauth"]
EXPOSE 8002
ADD zenauth /zenauth
