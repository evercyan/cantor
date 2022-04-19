import * as models from './models';

export interface go {
  "backend": {
    "App": {
		BatchUploadFile():Promise<models.Response>
		CheckFile(arg1:string):Promise<Error>
		CopyFileUrl(arg1:string):Promise<models.Response>
		DeleteFile(arg1:string):Promise<models.Response>
		GetConfig():Promise<models.Response>
		GetList():Promise<models.Response>
		Menu():Promise<models.Menu>
		OnBeforeClose(arg1:models.Context):Promise<boolean>
		SetConfig(arg1:string):Promise<models.Response>
		SyncDatabase():Promise<void>
		UpdateFileName(arg1:string,arg2:string):Promise<models.Response>
		Upload(arg1:string,arg2:Array<boolean>):Promise<string|Error>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
