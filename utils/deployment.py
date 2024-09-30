import boto3
import s3fs

def upload_model_version_to_s3(model_name,local_path):
    try:
        bucket_name = model_name
        
        s3_client = boto3.client('s3')
        response = s3_client.list_objects_v2(Bucket=bucket_name, Prefix=f"models")
        
        # Extract folder names and filter numerical ones
        folder_names = []
        for obj in response.get('Contents', []):
            match = re.search(rf"models/(\d+)/", obj['Key'])
            if match:
                folder_names.append(int(match.group(1)))
        
        new_version = max(folder_names, default=0) + 1
        
        object_key = f"models/{new_version}"
            
        s3_file = s3fs.S3FileSystem()
        s3_path = f"{bucket_name}/{object_key}"
        s3_file.put(local_path, s3_path, recursive=True) 
    
        return (s3_path,new_version)
    except Excetpion as e:
        print(f"Failed to upload model '{model_name}' to s3. {e}")
