import { Menu } from '@mantine/core';
import { useDidUpdate, useDisclosure } from '@mantine/hooks';
import { ChevronDown, ChevronUp } from 'lucide-react';
import type React from 'react';
import { Fragment, useCallback, useMemo, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router';

import classes from './DropdownSelect.module.css';

interface DropdownProps {
	defaultValue: string;
	defaultLabel: string;
	ariaLabel: string;
	records: Record<string, React.ReactNode>;

	icon?: React.ComponentType;
	customActions?: Record<string, () => void>;
	searchParamKey?: string;
	separatorValues?: string[];
	shouldUseNavigate?: boolean;
}

const DropdownSelect = ({
	defaultValue,
	defaultLabel,
	ariaLabel,
	records,
	icon: Icon,
	customActions,
	shouldUseNavigate = false,
	searchParamKey,
	separatorValues = [],
}: DropdownProps) => {
	const [searchParams, setSearchParams] = useSearchParams();
	const navigate = useNavigate();

	const option = searchParamKey
		? (searchParams.get(searchParamKey) ?? defaultValue)
		: defaultValue;

	const [radio, setRadio] = useState(option);
	const [open, { close, set: setOpen }] = useDisclosure(false);

	const handleOptionSubmit = useCallback(
		(value: string) => {
			if (shouldUseNavigate) {
				navigate(`/${value}`, { relative: 'route' });
			} else {
				setSearchParams((params) => {
					if (searchParamKey) {
						params.set(searchParamKey, value);
					}
					return params;
				});
			}
			setRadio(value);
		},
		[navigate, searchParamKey, setSearchParams, shouldUseNavigate],
	);

	const options = useMemo(
		() =>
			Object.entries(records).map(([value, content]) => {
				const option = (
					<Menu.Item
						key={value}
						className={classes.option}
						data-state={value === radio ? 'checked' : 'unchecked'}
						style={{
							color:
								value === radio ? 'var(--logo-green)' : 'var(--text-light)',
							fontWeight: value === radio ? 600 : undefined,
						}}
						onClick={() => {
							const customAction = customActions?.[value];
							if (customAction) {
								customAction();
								close();
								return;
							}

							handleOptionSubmit(value);
							close();
						}}
					>
						{content}
					</Menu.Item>
				);

				if (separatorValues.includes(value) && !shouldUseNavigate) {
					return (
						<Fragment key={value}>
							{option}
							<Menu.Divider className={classes.separator} />
						</Fragment>
					);
				}

				return option;
			}),
		[
			records,
			radio,
			customActions,
			separatorValues,
			shouldUseNavigate,
			close,
			handleOptionSubmit,
		],
	);

	useDidUpdate(() => {
		setRadio(option);
	}, [searchParams]);

	const labelComp =
		typeof records[option] === 'string' ? records[option] : defaultLabel;

	return (
		<Menu
			opened={open}
			onChange={setOpen}
			position="bottom-start"
			offset={8}
			withinPortal
			classNames={{ dropdown: classes.dropdown }}
		>
			<Menu.Target>
				<button
					type="button"
					className={classes.trigger}
					data-left={Boolean(Icon)}
					aria-label={ariaLabel}
				>
					<div>
						{Icon && <Icon />}
						<span>{labelComp}</span>
					</div>
					{open ? <ChevronUp size={15} /> : <ChevronDown size={15} />}
				</button>
			</Menu.Target>
			<Menu.Dropdown>{options}</Menu.Dropdown>
		</Menu>
	);
};

export { DropdownSelect };
