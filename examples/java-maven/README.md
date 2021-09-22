# Build Spec File Example using Java and Maven

This is an example Hello World project using JDK 15 and Maven. With the [OCI DevOps service](https://www.oracle.com/devops/devops-service/) and this project, you'll be able to build this application, package it as a JAR file and store it in [OCI Artifact Registry.](https://docs.oracle.com/en-us/iaas/artifacts/using/overview.htm)

In this example, you'll install JDK 15, build a JAR file of this Java Hello World Spring Boot Web app, and store your JAR file in the OCI Artifact Registry,  all using the OCI DevOps service!

Let's go!

## Building the application locally

### Download the repo
The first step to get started is to download the repository to your local workspace

```shell
git clone git@github.com:dlcbld/BuildSpecSampleFiles.git
cd java-maven
```

### Install and run the application

1. Install JDK 15 on your system: https://download.java.net/openjdk/jdk15/ri/openjdk-15+36_linux-x64_bin.tar.gz
2. Compile the app: 
   
   ```mvn clean install```
3. To verify, make sure the JAR file is created in your target folder.

## Building the application in OCI DevOps
Now that you've seen you can locally build this app, let's try this out through OCI DevOps Build service.

### Create External Connection to your Git repository 

1. Create a [DevOps Project](https://docs.oracle.com/en-us/iaas/Content/devops/using/devops_projects.htm) or use and an existing project. 
2. In your DevOps project, create an External Connection to your GitHub repository which holds your application.
   - Create a Personal Access Token (PAT): https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token
   - In the OCI Console, Go to Identity & Security -> Vault and create a [Vault]( https://docs.oracle.com/en-us/iaas/Content/KeyManagement/Concepts/keyoverview.htm) in compartment of your own choice.
   - Create a Master Key that will be used to encrypt the PATs. 
   - Select Secrets from under Resources and create a secret using PAT obtained from GitHub account.
   - Make a note of the OCID of the secret.
   - Now, go to the desired project and select External Connection from the resources.
   - Select type as GitHub and provide OCID of the secret under Personal Access Token.
   - Finally, allow BuildPipeline (dynamic group with DevOps Resources) to use a PAT secret by writing a policy in the root compartment as: ``` Allow dynamic-group dg-with-devops-resources to manage secret-family in tenancy```

### Setup your Build Pipeline
Create a new Build Pipeline to build, test and deliver artifacts from your GitHub Repository.

### Managed Build stage
In your Build Pipeline, first add a Managed Build stage
1. The Build Spec File Path is the relative location in your repo of the build_spec.yaml . Leave the default, for this example.
2. For the Primary Code Repository follow the below steps
    - Select connection type as GitHub
    - Select the external connection you created above
    - Give the HTTPS URL to the repo which contains your application.
    - Select main branch.
    
### Create a Artifact Registry repository
Reference: https://docs.oracle.com/en-us/iaas/artifacts/using/overview.htm
1. Go to Artifact Registry and click on Create Repository.
2. Fill in the details and click on Create and the repo will be created. Note down the ocid for the repo created.

<img src="./assets/Creating Registry.png" />

### Create a DevOps Artifact for your artifact repository
1. Next go to your project and in the Artifacts section, click on Add artifact.
2. Fill in the details as below:
    - Name: java-maven-artifact
    - Type: General Artifact
    - Artifact Source : Artifact Registry repository
    - Select Artifact registry repository as the one created in previous step
    - Artifact location: Choose set Custom location and provide path(target/HelloWorld-0.0.1-SNAPSHOT.jar) and version.
    - Replace parameters: Yes, substitute placeholders
   
<img src="./assets/Creating DevOps Artifact.png" />

Required policies must be added in the root compartment for the Artifact Registry repository and DevOps Artifact resource.
1. Provide access to Artifact Repository to deliver artifacts : ```Allow dynamic-group dg-with-devops-resources to manage repos in tenancy```
2. Provide access to read deploy artifacts in deliver artifact stage : ```Allow dynamic-group dg-with-devops-resources to manage devops-family in tenancy```


### Add a Deliver Artifacts stage
Let's add a **Deliver Artifacts** stage to your Build Pipeline to deliver the JAR file to an artifact repository.

The Deliver Artifacts stage **maps** the ouput artifacts from the Managed Build stage with the version to deliver to the artifact repository, through the DevOps Artifact resource.

Add a **Deliver Artifacts** stage to your Build Pipeline after the **Managed Build** stage. To configure this stage:
1. In your Deliver Artifacts stage, choose `Select Artifact`
   <img src="./assets/Final Details in Deliver Artifact Stage.png" />
1. From the list of artifacts select the `java-maven-artifact` artifact that you created above
   <img src="./assets/Selecting DevOps Artifact resource.png" />
1. In the next section, you'll assign the  container image outputArtifact from the `build_spec.yaml` to the DevOps project artifact. For the "Build config/result Artifact name" enter: `HelloWorld-Jar`
   <img src="./assets/Final Details in Deliver Artifact Stage.png" />


### Run your Build in OCI DevOps

#### From your Build Pipeline, choose `Manual Run`
Use the Manual Run button to start a Build Run

Manual Run will use the Primary Code Repository, will start the Build Pipeline, first running the Managed Build stage, followed by the Deliver Artifacts stage.

After the Build Pipeline execution is complete, we can view the JAR file stored in the Artifact Registry, which can then be pulled to local workspace (Click on the ellipses for more options and choose ``` Download```).

