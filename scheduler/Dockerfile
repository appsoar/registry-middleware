FROM centos:latest


EXPOSE 9090
RUN mkdir /host
RUN mkdir /host/proc
VOLUME ["/host/proc"]
COPY scheduler /scheduler/
RUN chmod +x /scheduler/scheduler

CMD ["/scheduler/scheduler"]
