

export interface UploadMetaData {
    uploadSpeeds: Record<string, number>,
    fileProgress: Record<string, number>
    averageSpeed: Record<string, number>,
    speedCounter: Record<string, number>,
    etaInMinutes: Record<string, number>
}
