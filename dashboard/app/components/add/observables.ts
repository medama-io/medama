import { observable } from '@legendapp/state';

const name$ = observable('');
const hostname$ = observable('');

export { hostname$, name$ };
