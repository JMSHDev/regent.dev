/**
* JetBrains Space Automation
* This Kotlin-script file lets you automate build activities
* For more info, see https://www.jetbrains.com/help/space/automation.html
*/


job("Compile agent and ExampleApp") {
    container("golang:buster") {
        shellScript {
            interpreter = "/bin/bash"
            location = "build_agent.sh"
        }
    }
}
