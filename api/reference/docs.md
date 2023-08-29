
**Reference Models and Populating Data**

The reference models in the application are critical for its proper functioning. These models store data that is essential for various operations. To ensure that these models are populated with the necessary data, you can utilize the `AbstractBaseReference` model. Additionally, each reference model should implement the `populate` method to handle the logic for populating data.

**Using the `smy` CLI Tool**

You can use the `smy` CLI tool to easily populate these reference models. To do this, follow these steps:

1. Open the `/Smylet/symlet-backend/utilities/cli/populate/populate.go` file.

2. Locate the `referenceModelMap` declaration in the file. This map associates each model's name with its corresponding struct. For example:

   ```go
   referenceModelMap := map[string]reference.ReferenceModelInterface{
       "amenities": reference.ReferenceHostelAmmenities{},
       "university": reference.ReferenceUniversity{},
   }
   ```

3. To add a new model to the reference model map, simply include the model's struct and name as a new entry in the map. For instance, if you have a new model named `newmodel`, add the following line:

   ```go
   "newmodel": reference.NewModelStruct{},
   ```

4. Update the AutoMigrate list in the same file by adding the new model's struct to it. This ensures that the model's table is created in the database. For example:

   ```go
   db.AutoMigrate(&reference.ReferenceHostelAmmenities{}, &reference.ReferenceUniversity{}, &reference.NewModelStruct{})
   ```

**Note:** Make sure to replace `NewModelStruct` with the actual name of the new model's struct.

By following these steps, you can seamlessly populate reference models and ensure that your application functions smoothly with the required data.

