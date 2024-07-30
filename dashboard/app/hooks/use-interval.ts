import { useEffect, useRef, useState } from 'react';

interface UseIntervalOptions {
	/** If set, the interval will start automatically when the component is mounted, `false` by default */
	autoInvoke?: boolean;
}

const useInterval = (
	fn: () => void,
	interval: number,
	{ autoInvoke = false }: UseIntervalOptions = {},
) => {
	const [active, setActive] = useState(false);
	const intervalRef = useRef<number>();
	const fnRef = useRef<() => void>();

	const start = () => {
		setActive((old) => {
			if (!old && !intervalRef.current) {
				// biome-ignore lint/style/noNonNullAssertion: guaranteed to be defined
				intervalRef.current = window.setInterval(fnRef.current!, interval);
			}
			return true;
		});
	};

	const stop = () => {
		setActive(false);
		window.clearInterval(intervalRef.current);
		intervalRef.current = undefined;
	};

	const toggle = () => {
		if (active) {
			stop();
		} else {
			start();
		}
	};

	// biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
	useEffect(() => {
		fnRef.current = fn;
		active && start();
		return stop;
	}, [fn, active, interval]);

	// biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
	useEffect(() => {
		if (autoInvoke) {
			start();
		}
	}, []);

	return { start, stop, toggle, active };
};

export { useInterval };
