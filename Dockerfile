FROM scratch
CMD ["/zenauth"]
EXPOSE 5000
EXPOSE 5001
ADD zenauth /zenauth
