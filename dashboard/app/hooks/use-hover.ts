import { useCallback, useEffect, useRef, useState } from 'react';

const useHover = <T extends HTMLElement = HTMLDivElement>() => {
	const [hovered, setHovered] = useState(false);
	const ref = useRef<T>(null);
	const onMouseEnter = useCallback(() => setHovered(true), []);
	const onMouseLeave = useCallback(() => setHovered(false), []);

	// biome-ignore lint/correctness/useExhaustiveDependencies: ref.current is not a dependency
	useEffect(() => {
		if (ref.current) {
			ref.current.addEventListener('mouseenter', onMouseEnter);
			ref.current.addEventListener('mouseleave', onMouseLeave);

			return () => {
				ref.current?.removeEventListener('mouseenter', onMouseEnter);
				ref.current?.removeEventListener('mouseleave', onMouseLeave);
			};
		}

		return undefined;
	}, []);

	return { ref, hovered };
};

export { useHover };
