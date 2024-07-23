import {
	type DependencyList,
	type EffectCallback,
	useEffect,
	useRef,
} from 'react';

/**
 * A hook that runs an effect only after the first render.
 */
export const useDidUpdate = (
	fn: EffectCallback,
	dependencies?: DependencyList,
) => {
	const mounted = useRef(false);

	useEffect(
		() => () => {
			mounted.current = false;
		},
		[],
	);

	// biome-ignore lint/correctness/useExhaustiveDependencies: No need to add fn to dependencies
	useEffect(() => {
		if (mounted.current) {
			return fn();
		}

		mounted.current = true;
		return undefined;
	}, dependencies);
};
