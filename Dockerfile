FROM scratch
ADD goshortest /

EXPOSE 8000

CMD ["/goshortest"]
