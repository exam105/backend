# Welcome to the Exam105 Backend Application


## Mission Statement: 
Exam105.com is a platform where the students and teachers will be able to get, search, and interact with the past papers of Cambridge Board Education System. It would provide its users with a lot of benefits. The users will be able to pick and choose from a wide range of papers conducted over the years, all the way from 2011, and make a customized pdf book to download and work with it online however they may prefer; be it to prepare for the upcoming exams (from students' perspective), or to make new papers (from teachers' perspective). And as we progress, there will be more functionalities provided by the system in future.

---

This is the backend application of Exam105. It uses Golang as the primary language. This application uses Eco framework, and it is coded using the clean architecture. The application uses docker to containerize and deploy using ansible. 

### Clean Architecture:
This application follows the clean architecture approach. Clean architecture helps get rid of many constraints. The architecture does not depend on the existence of some library or framework. It is easily testable; the business rules can be tested without depending upon the UI, Database, Web Server, or any other external element. The business rules are not bound to the database, that means we can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or any other database. Below is the diagram of clean archiecture that is implemented in this application. 

![CleanArchitecture](https://user-images.githubusercontent.com/52341921/133215704-d766651e-7e58-4235-80dd-6293ae09c115.jpg) 

*Credit: [Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)* 


---

## System Architecture:
Now we describe about the system architecture of Exam105.com as a whole. The system architecture is based on AWS technologies. The overall system is kept very cost effective. Below is the diagram of system architecture.

![exam105_system_architecture](https://user-images.githubusercontent.com/52341921/133076841-609da4af-87bc-4bef-9867-489850b5ea88.png)

The technologies that our system uses are listed down below.
### 1. EC2 Instances: 
   Amazon Elastic Compute Cloud (Amazon EC2) provides scalable computing capacity in the Amazon Web Services (AWS) Cloud. Using Amazon EC2 eliminates the need to invest in hardware up front, so we can develop and deploy applications faster. Thus, we have used Amazon EC2 to launch three virtual servers. It allows us to configure security and networking, and manage storage.
### 2. VPC (Virtual Private Cloud):
   A virtual private cloud (VPC) is a virtual network dedicated to your AWS account. It is logically isolated from other virtual networks in 	the AWS Cloud. You can launch your AWS resources, such as Amazon EC2 instances, into your VPC. You can specify an IP address range for 		the VPC, add subnets, associate security groups, and configure route tables.
### 3. S3:
   Amazon S3 or Amazon Simple Storage Service is a service offered by Amazon Web Services that provides object storage through a web service 	interface. We have used this service to store the images of our system.
### 4. Route53:
   Amazon Route 53 is a scalable and highly available Domain Name System service. We have utilized Route53 to provide the domain name for our system.
### 5. Lightsail:
   Lightsail provides developers compute, storage, and networking capacity and capabilities to deploy and manage websites and web applications in the cloud. We have used the lightsail to test and deploy the application. We have 2 machines on lightsail. 1 machine is for development and other for deployment to dev and production.
### 6. Docker:
   Docker is an operating system for containers. Containers virtualize the operating system of a server. Docker is installed on each server	of the system architecture and provides simple commands you can use to build, start, or stop containers.
### 7. Linux (Ubuntu):
   Ubuntu is installed on all of the nodes of the architecture. Linux provides a secure, stable, and high performance execution environment to develop and run cloud and enterprise applications.

### Servers/Nodes:
   The system is composed of three servers in total. All the three servers are located at Mumbai region, in three geographically separated availability zones. This is done so that, if in case some unexpected damage occurs at one area of the availability zones, the rest available servers will handle all the jobs and this will not halt the system from processing its normal operations. Linux (Ubuntu) OS, nginx and docker are installed on each of the nodes.

### Proxy server: 
   We have Nginx docker image acting as reverse proxy. This node is known as Manager, while the other two nodes are known as Worker. Because SWARM Manager is installed here.  This node is running nginx docker image. The swarm orchestration performs load balancing, meaning, if one machine goes down, the swarm will run the operations on the other available machines and not let the system go down. 
When a user request arrives, it first goes to the proxy server, and from there it is forwarded to the respective servers. Round Robin is the formula that is applied on the proxy server, that is, first request goes to the first server, and the second request to the second server.

The servers other than proxy server, have exam105_frontend and exam105_backend docker images.

### Databases: 
   Mongo DB is used for the storing of system data. The db is in replicated mode. Meaning the same db instance is available at multiple places. And one of the db instances is set as master, while the rest are known as slaves. All the replicas are connected to eachother. The databases can be connected via private network. 

### Swarm Orchestration:
   The swarm orchestration performs load balancing, meaning, if one machine goes down, the swarm will run the operations on the other available machines and not let the system go down. The swarmpit is a web UI for Docker Swarm, which provides the information about all of the machines, the services that are working, or not working, all of the stats of the machines that the architecture is made up of.

### System description:
   1 cluster (cluster is the network of machines.)
Total machines are 3. 6 CPUS (2 cpus per machine)
Services: exam105_nginx (load balancer), exam105_backend, exam105_frontend

### Performance:
   We performed the load testing of our web application using (JMETER ICON)JMeter. JMeter is an Apache project that can be used as a load testing tool for analyzing and measuring the performance of web applications.
#### Testing parameters:
   Number of Threads (users): The number of users that JMeter would attempt to simulate were set to 1000.
   Ramp-Up Period (in seconds): The duration of time that JMeter would distribute the start of the threads over was set to 10 seconds.
   Loop Count: The number of times to execute the test was set to 1.
   
### Results:
The results are shown in the screenshots below (these are some results, not all as we could not capture all due to the large number of data). According to the **View Results in Table** output, the **Status** of all the requests was “Success” (indicated by a green shield with a tick in it); the range of both the **Sample Time** and **Latency** was 179-2601 ms.

![Screenshot (13)](https://user-images.githubusercontent.com/52341921/133097602-54cac6e5-25a0-47c8-ad2d-203b024d9992.png)

![Screenshot (12)](https://user-images.githubusercontent.com/52341921/133097645-c3265c7d-d3e3-4594-ac87-8bdb6212855d.png)
