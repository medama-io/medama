import * as DropdownMenu from '@radix-ui/react-dropdown-menu';
import { useNavigate, useSearchParams } from '@remix-run/react';
import { ChevronDown, ChevronUp } from 'lucide-react';
import type React from 'react';
import { Fragment, useMemo, useState } from 'react';

import { useDidUpdate } from '@/hooks/use-did-update';

import classes from './DropdownSelect.module.css';

interface DropdownProps {
	defaultValue: string;
	defaultLabel: string;
	ariaLabel: string;
	records: Record<string, string>;

	icon?: React.ComponentType;
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
	shouldUseNavigate = false,
	searchParamKey,
	separatorValues = [],
}: DropdownProps) => {
	const [searchParams, setSearchParams] = useSearchParams();
	const navigate = useNavigate();

	const option = searchParamKey
		? searchParams.get(searchParamKey) ?? defaultValue
		: defaultValue;

	const [radio, setRadio] = useState(option);
	const [open, setOpen] = useState(false);

	const options = useMemo(
		() =>
			Object.entries(records).map(([value, label]) => {
				const option = (
					<DropdownMenu.RadioItem
						key={value}
						className={classes.option}
						value={value}
					>
						{label}
					</DropdownMenu.RadioItem>
				);

				if (separatorValues.includes(value) && !shouldUseNavigate) {
					return (
						<Fragment key={value}>
							{option}
							<DropdownMenu.Separator className={classes.separator} />
						</Fragment>
					);
				}

				return option;
			}),
		[records, separatorValues, shouldUseNavigate],
	);

	const handleOptionSubmit = (value: string) => {
		if (shouldUseNavigate) {
			navigate(`/${value}`, { relative: 'route' });
		} else {
			setSearchParams((prevParams) => {
				const newParams = new URLSearchParams(prevParams);
				if (searchParamKey) {
					newParams.set(searchParamKey, value);
				}
				return newParams;
			});
		}
		setRadio(value);
	};

	useDidUpdate(() => {
		setRadio(option);
	}, [searchParams]);

	return (
		<DropdownMenu.Root onOpenChange={(isOpen) => setOpen(isOpen)} modal={false}>
			<DropdownMenu.Trigger asChild>
				<button
					type="button"
					className={classes.trigger}
					data-left={Boolean(Icon)}
					aria-label={ariaLabel}
				>
					<div>
						{Icon && <Icon />}
						<span>{records[option] ?? defaultLabel}</span>
					</div>
					{open ? <ChevronUp /> : <ChevronDown />}
				</button>
			</DropdownMenu.Trigger>
			<DropdownMenu.Portal>
				<DropdownMenu.Content className={classes.dropdown} sideOffset={8}>
					<DropdownMenu.RadioGroup
						value={radio}
						onValueChange={handleOptionSubmit}
					>
						{options}
					</DropdownMenu.RadioGroup>
				</DropdownMenu.Content>
			</DropdownMenu.Portal>
		</DropdownMenu.Root>
	);
};

export { DropdownSelect };
