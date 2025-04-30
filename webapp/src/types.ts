
export interface File {
    id: number;
    Filepath: string;
    Uploaded: boolean;
    UploadedAt?: Date;
    CreatedAt: Date;
    UpdatedAt: Date;
}

interface FileJson {
    id: number,
    CreatedAt: string,
    UpdatedAt: string,
    UploadedAt?: string,
    Filepath: string;
    Uploaded: boolean
}

interface UploadProcessJson {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string,
    Files: FileJson[],
}

export interface UploadProcess {
    ID: number;
    CreatedAt: Date,
    UpdatedAt: Date,
    Files: File[];
}

export function mapUploadProcess(json: UploadProcessJson): UploadProcess {
    return {
        ...json,
        CreatedAt: new Date(json.CreatedAt),
        UpdatedAt: new Date(json.UpdatedAt),
        Files: json.Files.map(mapFile)
    }
}

export function mapFile(json: FileJson): File {
   return {
       ...json,
       UpdatedAt: new Date(json.UpdatedAt),
       CreatedAt: new Date(json.CreatedAt),
       UploadedAt: json.UploadedAt ? new Date(json.UploadedAt) : undefined
   }
}