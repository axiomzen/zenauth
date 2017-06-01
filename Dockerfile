FROM scratch
CMD ["/zenauth"]
EXPOSE 5000
ADD zenauth /zenauth
