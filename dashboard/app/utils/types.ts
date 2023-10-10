export interface GetUser {
	email: string;
	language: string;
	dateUpdated: number;
	dateCreated: number;
}

export interface PostUser {
	email: string;
	password: string;
	language?: string;
}
