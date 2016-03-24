FROM centos:latest

COPY scheduler /scheduler/
RUN chmod +x /scheduler/scheduler

EXPOSE 9090
RUN mkdir /host
RUN mkdir /host/proc
VOLUME ["/host/proc"]

CMD ["/scheduler/scheduler"]
