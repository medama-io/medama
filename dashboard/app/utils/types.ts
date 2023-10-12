export const SESSION_NAME = '_me_sess';

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
