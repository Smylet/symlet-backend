**Using the `smy` CLI Tool**

The `smy` CLI tool is designed to simplify common tasks, such as populating reference tables. To extend this tool with new commands, you can follow these steps:

1. Create a New Command:

   To add a new command to the `smy` CLI tool, start by creating a new folder to house the command code. Inside this folder, create a new Go file, e.g., `newcommand.go`, and implement the functionality of your new command.

2. Implement the Command Code:

   In the `newcommand.go` file, define the logic for your new command using the Cobra framework. This includes setting up flags, arguments, and specifying the action to be taken when the command is executed.

3. Add the Command to the Root Command:

   In the `root.go` file, which is located in the `root` package, you'll find an `init` function. Inside this function, add the new command to the root command using the `AddCommand` method. For example:

   ```go
   // Inside root/root.go
   func init() {
       // Other code...
       rootCmd.AddCommand(newcommandpkg.NewCommand)
   }
   ```

   Replace `newcommandpkg.NewCommand` with the actual reference to your new command's implementation.

4. Usage:

   After completing these steps, you can use the `smy` CLI tool to run your new command. The syntax will be similar to:

   ```
   smy newcommand [flags]
   ```

**Note:** Ensure that you replace `newcommand` with the actual name of your new command and follow the correct flags and arguments syntax.