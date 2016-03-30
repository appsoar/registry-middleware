FROM centos:latest


EXPOSE 9090
RUN mkdir /host
RUN mkdir /host/proc
RUN mkdir /.hidden
RUN mkdir /.hidden/root
VOLUME ["/.hidden/root"]
VOLUME ["/host/proc"]
COPY scheduler /scheduler/
RUN chmod +x /scheduler/scheduler

CMD ["/scheduler/scheduler"]
