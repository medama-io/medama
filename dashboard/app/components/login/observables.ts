import { observable } from '@legendapp/state';

const username$ = observable('');
const password$ = observable('');

export { password$, username$ };
