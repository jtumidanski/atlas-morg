FROM maven:3.6.3-openjdk-14-slim AS build

COPY settings.xml /usr/share/maven/conf/

COPY pom.xml pom.xml
COPY morg-api/pom.xml morg-api/pom.xml
COPY morg-model/pom.xml morg-model/pom.xml
COPY morg-base/pom.xml morg-base/pom.xml

RUN mvn dependency:go-offline package -B

COPY morg-api/src morg-api/src
COPY morg-model/src morg-model/src
COPY morg-base/src morg-base/src

RUN mvn install -Prunnable

FROM openjdk:14-ea-jdk-alpine
USER root

RUN mkdir service

COPY --from=build /morg-base/target/ /service/

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait

RUN chmod +x /wait

ENV JAVA_TOOL_OPTIONS -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005

EXPOSE 5005

CMD /wait && java --enable-preview -jar /service/morg-base-1.0-SNAPSHOT.jar -Xdebug