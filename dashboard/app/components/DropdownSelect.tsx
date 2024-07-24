import * as DropdownMenu from '@radix-ui/react-dropdown-menu';
import { useNavigate, useSearchParams } from '@remix-run/react';
import { ChevronDown, ChevronUp } from 'lucide-react';
import type React from 'react';
import { Fragment, useMemo, useState } from 'react';

import { useDidUpdate } from '@/hooks/use-did-update';
import { useDisclosure } from '@/hooks/use-disclosure';

import classes from './DropdownSelect.module.css';

interface DropdownProps {
	defaultValue: string;
	defaultLabel: string;
	ariaLabel: string;
	records: Record<string, React.ReactNode>;

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
	const [open, { close, toggle }] = useDisclosure(false);

	const options = useMemo(
		() =>
			Object.entries(records).map(([value, content]) => {
				const option =
					typeof content === 'string' ? (
						<DropdownMenu.RadioItem
							key={value}
							className={classes.option}
							value={value}
						>
							{content}
						</DropdownMenu.RadioItem>
					) : (
						<DropdownMenu.Item key={value} asChild onSelect={close}>
							{content}
						</DropdownMenu.Item>
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
		[records, separatorValues, shouldUseNavigate, close],
	);

	const handleOptionSubmit = (value: string) => {
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
	};

	useDidUpdate(() => {
		setRadio(option);
	}, [searchParams]);

	const labelComp =
		typeof records[option] === 'string' ? records[option] : defaultLabel;

	return (
		<DropdownMenu.Root open={open} onOpenChange={toggle} modal={false}>
			<DropdownMenu.Trigger asChild>
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
