## Cloudflow - AWS API Request Logger

```
	 ██████╗██╗      ██████╗ ██╗   ██╗██████╗ ███████╗██╗      ██████╗ ██╗    ██╗
	██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗██╔════╝██║     ██╔═══██╗██║    ██║
	██║     ██║     ██║   ██║██║   ██║██║  ██║█████╗  ██║     ██║   ██║██║ █╗ ██║
	██║     ██║     ██║   ██║██║   ██║██║  ██║██╔══╝  ██║     ██║   ██║██║███╗██║
	╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝██║     ███████╗╚██████╔╝╚███╔███╔╝
	 ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝ 
```																			 		
Cloudflow is a simple tool that listens for and logs AWS API requests made using AWS Client-Side Monitoring (CSM). This allows you to analyze and debug API interactions locally.
<img src=cloudflow.png>
### Features

* Captures AWS API request data sent via CSM.
* Outputs data in JSON format to a user-defined file.
* Optionally displays data to the console in raw or formatted table view.
* Verifies if AWS CSM is enabled on the system.

### Prerequisites

In order for cloudflow to receive log events the following environment variables must be set in the shell that will be making the AWS API calls.
```
export AWS_CSM_ENABLED=true
export AWS_CSM_PORT=31000
export AWS_CSM_HOST=127.0.0.1
```

### Installation

**Clone the Cloudflow repository:**

   Use `git` to clone the Cloudflow repository from GitHub:

   ```bash
   git clone https://github.com/louisbarrett/cloudflow
   ```

 **Build Cloudflow:**

   Navigate to the cloned directory and build Cloudflow using Go:

   ```bash
   cd cloudflow
   go build -o cloudflow .
   ```

### Usage

**Run Cloudflow:**

   Run Cloudflow with the desired options:

   - `-port <port number>` (default: 31000): Port to listen for UDP traffic on.
   - `-output <file name>` (default: output.jsonl): File to write captured data.
   - `-verbose`: Print raw data to the console.
   - `-pretty`: Display data in a formatted table on the console.
   - `-doctor`: Checks if AWS CSM environment variables are set.
   - `-s`: Do not write to standard output

   **Warning** 
   When `-pretty` is not specified SessionTokens will be sent to standard output.


   **Example:**

   ```bash
   ./cloudflow -port 31000 -output api_requests.log -pretty
   ```

   **Sample Output**
   ```
   Timestamp      AccessKey             Service  Api          Region     UserAgent
1713392851294  ASIA1337EXAMPLEDMP  S3       ListBuckets  us-west-2  aws-cli/2.15.35 Python/3.11.8 Darwin/23.4.0 exe/x86_64 prompt/off command/s3.ls
1713392862041  ASIA1337EXAMPLEDMP  IAM      ListUsers    us-east-1  aws-cli/2.15.35 Python/3.11.8 Darwin/23.4.0 exe/x86_64 prompt/off command/iam.list-users
1713392872909  ASIA1337EXAMPLEDMP  KMS      ListKeys     us-west-2  aws-cli/2.15.35 Python/3.11.8 Darwin/23.4.0 exe/x86_64 prompt/off command/kms.list-keys

   ```

### Use Cases

* **Adversarial Simulation Testing:** Use Cloudflow to capture and analyze API requests made during security testing that simulates attacks against your AWS environment.
* **Debugging Services:** Identify and debug issues in your services that interact with the AWS API by monitoring the specific requests being made.
