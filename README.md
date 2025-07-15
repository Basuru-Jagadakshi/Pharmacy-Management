# GO + WSO2 API Manager Integration

In this tutorial, I'm going to show how WSO2 API manager is integrated with the my GO backend and expose my GO backend via WSO2 API gateway.

## Step 1 :- Run the GO backend and PostgreSQL using docker
Here, 
    - GO backend(port - 8080)
    - PostgreSQL (port - 5432)

Then run below command,

docker compose up --build

After up all the containers, check local endpoints by using postman via below endpoints.

    -GET    http://localhost:8080/drugs
    -POST   http://localhost:8080/drugs
    -PUT    http://localhost:8080/drugs/{id}
    -DELETE http://localhost:8080/drugs/{id}

Then type,

docker compose down

it downs all the container.




## Step 2 :- WSO2 API manager 

Add below code into docker compose file to get WSO2 API manager 4.0.0 image.

    wso2apim:
        image: wso2/wso2am:4.0.0
        container_name: wso2_api_manager
        ports:
        - "9443:9443"   
        - "8243:8243"   
        depends_on:
        - app
        restart: unless-stopped

And then agin up all the containers bu using,
    
    docker compose up --build

![image alt](https://github.com/Basuru-Jagadakshi/Pharmacy-Management/blob/31d499287372a06db4cb4b168578750454bc4a2b/Screenshot%202025-07-14%20at%2003.45.41.png)



## Step 3 :- Logging into WSO2 publisher portal

In your web browser search below url,
    https://localhost:9443/publisher

Then give default crdentials as,
    username = admin
    password = admin


![image alt](https://github.com/Basuru-Jagadakshi/Pharmacy-Management/blob/4498938eb71edb14defda2640a5e3f4f5b42e3f6/Screenshot%202025-07-14%20at%2003.46.46.png)


## Step 4 :- Create REST API in WSO2 API manager

After logging you can see the admin panel in publisher portal. Click on "REST API" tab and select "Start From Scratch". Now You will have a form. Fille it like below,

    Name - GoDrugAPI
    Context - /drugsapi
    Version - 1.0.0
    Endpoint - http://host.docker.internal:8080

Then click "Create" button.



## Step 5 :- Add Resources

Then, You will see your GoDrugAPI panel. Then go to "Resource" tab at the left corner to add resources. And then add those like below,

    GET /drugs
    POST    /drugs
    PUT /drugs/{id}
    DELETE  /drugs/{id}

After adding all the resources click "Save" button.



## Step 6 :- Depoly and Publish

Then go to "Deployments" tab and click "Deploy to Gateway" and go to "Lifecycle" tab. There, you will see API status is created. Click "Publish" button. In that case, API status will be changed as "Published". 

Note :- This step is must. Otherwise you won't be able to access your API in the developer portal.

Now you can your access your endpoint like below examples,
    https://localhost:8243/drugsapi/1.0.0/drugs

But you need access token. Let's see how it is done.




## Step 7 :- Test via postman

Then go to Developer via below link,
    https://localhost:9443/devportal

    username = admin
    password = admin

Now, you will see you GoDrugAPI in the panel. Click on it and subscibe it. Through DefaulApplication generate access token check your API. Otherwise you can check your APIs developer portal as well. For that, Go to "Try Out" section.






## Conclusion

This project demonstrates how to integrate a RESTful Go backend with WSO2 API Manager using Docker.

Key points:
- Defined and published secure APIs using WSO2 Publisher
- Subscribed and tested them using Developer Portal and Postman
- Used Docker Compose to run PostgreSQL, Go API, and API Manager

This setup can be extended to any real world microservice that needs authentication, throttling, and centralized API management.
