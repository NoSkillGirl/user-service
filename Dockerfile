FROM ubuntu:16.04
COPY user-service /user-service
RUN chmod +x /user-service
CMD /user-service
EXPOSE 8082