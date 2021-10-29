# Build Spec File Example using Java and Docker

This is an example Hello World project using Java SE JDK 11 and Docker. With the 
[OCI DevOps service](https://www.oracle.com/devops/devops-service/) and this project, 
you'll be able to build this application, build a Docker image and store it in 
[DockerHub](https://hub.docker.com/)

In this example, you'll build a container image of this Java Hello World app, and push your built image 
to your DockerHub repository, all using the OCI DevOps service!

Let's go!

## Building the application locally

### Download the repo
The first step to get started is to download the repository to your local workspace

```shell
git clone git@github.com:dlcbld/BuildSpecSampleFiles.git
cd examples/java-dockerhub
```

### Install and run the application

1. Install Java SE SDK 11 on your system: https://java.com/en/download/help/download_options.html
2. Compile the app:

   ```javac src/com/sample/HelloWorld.java```
3. Run the app:

   ```java src/com/sample/HelloWorld```
4. To verify, make sure the string "Hello World" is printed in the shell.

### Build a container image for the app
You can locally build a container image using docker (or your favorite container image builder), to verify that you can run the app within a container.

```
docker build -t hello-world:1.0 .
```

Verify that your image was built, with `docker images`

Next run your local container and confirm you can access the app running in the container
```
docker run --rm -d --name hello-world:latest
```

The string "Hello World" must be printed in the shell if your container has been built successfully.

## Building the application in OCI DevOps
Now that you've seen you can locally build this app, let's try this out through OCI DevOps Build service.

   
### Creating DockerHub and GitHub Secrets

1. In the OCI Console, Go to Identity & Security -> Vault and create a [Vault]( https://docs.oracle.com/en-us/iaas/Content/KeyManagement/Concepts/keyoverview.htm) in compartment of your own choice.
2. Create a Master Key that will be used to encrypt the Secrets(login credentials).
3. Select Secrets from under Resources and create one each for the following:
   - Personal Access Token(PAT) obtained from GitHub account (Check [this](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token) to see how you can make one)
   - DockerHub Username
   - DockerHub Password
4. Make a note of the OCID of the secret.
5. Finally, allow BuildPipeline (dynamic group with DevOps Resources) to use secrets by writing a policy in the root compartment as:``` Allow dynamic-group dg-with-devops-resources to manage secret-family in tenancy```


### Create External Connection to your Git repository

1. Create a [DevOps Project](https://docs.oracle.com/en-us/iaas/Content/devops/using/devops_projects.htm) or open an existing project.
2. Select External Connection from the resources. 
3. Select type as GitHub and provide OCID of the secret holding your Personal Access Token.

### Setup your Build Pipeline
Create a new Build Pipeline to build, test and deliver artifacts from your GitHub Repository.

### Managed Build stage
In your Build Pipeline, add a Managed Build stage
1. The Build Spec File Path is the relative location in your repo of the build_spec.yaml . Leave the default, for this example.
2. For the Primary Code Repository follow the below steps
   - Select connection type as GitHub
   - Select the external connection you created above
   - Give the HTTPS URL to the repo which contains your application.
   - Select main branch.


### Run your Build in OCI DevOps

#### From your Build Pipeline, choose `Manual Run`
Use the Manual Run button to start a Build Run

Manual Run will use the Primary Code Repository, will start the Build Pipeline, running the Managed Build stage.

After the Build Pipeline execution is complete, we can view the container image stored in your DockerHub repository.


