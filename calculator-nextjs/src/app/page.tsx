"use client";
import Image from "next/image";
import styles from "./page.module.css";
import {useState} from "react";
import {client} from "@/rpc/client";

export default function Home() {
    // 表达式
    const [expression, setExpression] = useState("0");
    // 当前值
    const [currentValue, setCurrentValue] = useState('0');
    // 错误信息
    const [error, setError] = useState("");
    // 点击等于判断
    const [isEqualsClicked, setIsEqualsClicked] = useState(false);

    // 等待操作数
    const [waitingForOperand, setWaitingForOperand] = useState(false);

    const [operation, setOperation] = useState<string | null>(null);


    // 显示结果渲染
    const renderDisplay = () => {
        const displayContent = error ? (
            <div className={styles.errorDisplay}>{error}</div>
        ) : (
            currentValue !== '0' ? `= ${currentValue}` : '0'
        );
        return (
            <>
                <div className={isEqualsClicked ? styles.display : styles.expressionDisplay}>
                    {expression || "0"}
                </div>
                <div className={isEqualsClicked ? styles.expressionDisplay : `${styles.display} ${styles.displayContentHidden}`}>
                    {displayContent}
                </div>
            </>
        );
    };

    // 切换正负号
    const handlePlusMinus = () => {
        const newValue = parseFloat(currentValue) * -1;
        setCurrentValue(String(newValue));
        setExpression(prev => {
            if (prev === currentValue) return String(newValue);
            return prev.replace(currentValue, String(newValue));
        });
    };

    // 百分比
    const handlePercentage = () => {
        const newValue = parseFloat(currentValue) / 100;
        setCurrentValue(String(newValue));
        setExpression(prev => {
            if (prev === currentValue) return String(newValue);
            return prev.replace(currentValue, String(newValue));
        });
    };

    // 处理运算符
    const performOperation = (nextOperation: string) => {
        if (isEqualsClicked) {
            setIsEqualsClicked(false);
        }

        // 如果当前正在等待操作数，只更新运算符
        if (waitingForOperand) {
            setExpression(prev => {
                // 替换最后一个运算符
                return prev.replace(/[\+\-\*\/]$/, nextOperation);
            });
            setOperation(nextOperation);
            return;
        }

        // 更新表达式
        setExpression(prev => {
            // 如果当前表达式是"0"或者是新开始的计算
            if (prev === "0" || isEqualsClicked) {
                return currentValue + nextOperation;
            }
            return prev + nextOperation;
        });

        setWaitingForOperand(true);
        setOperation(nextOperation);
    };

    // 输入数字
    const inputDigit = (digit: string) => {
        if (isEqualsClicked) {
            clearAll();
        }

        if (waitingForOperand) {
            setCurrentValue(digit);
            setWaitingForOperand(false);
        } else {
            setCurrentValue(prev => prev === '0' ? digit : prev + digit);
        }

        // 更新表达式
        setExpression(prev => {
            if (prev === "0" || isEqualsClicked) {
                return digit;
            }
            return prev + digit;
        });
    };

    // 输入小数点
    const inputDot = () => {
        if (isEqualsClicked) {
            clearAll();
            setCurrentValue('0.');
            setExpression('0.');
            return;
        }

        if (waitingForOperand) {
            setCurrentValue('0.');
            setWaitingForOperand(false);
            setExpression(prev => prev + '0.');
            return;
        }

        if (!currentValue.includes('.')) {
            setCurrentValue(prev => prev + '.');
            setExpression(prev => prev + '.');
        }
    };

    // AC
    const clearAll = () => {
        setIsEqualsClicked(false)
        setCurrentValue('0');
        setExpression('0');
        setOperation(null);
        setWaitingForOperand(false);
        setError("");

    };

    const calculateResult = () => {
        if (waitingForOperand || !operation) return;
        setIsEqualsClicked(true);
        try {
            // 调用gRPC后端服务
            const response = client.calculate({expression});
            setCurrentValue(response.result);
        } catch (err) {
            // 提取错误消息中]后面的部分
            const errorMessage = err.message.match(/](.*)$/);
            if (errorMessage) {
                setError(errorMessage[1].trim());
            } else {
                setError("未知错误");
            }
            setCurrentValue('Error');
        }
    };


    return (
        <div className={styles.page}>
            <main className={styles.main}>
                <div className={styles.title}>简易计算器</div>

                <div className={styles.calculator}>
                    {renderDisplay()}

                    <div className={`${styles.buttonRow} ${styles.firstBtn}`}>
                        <button className={styles.functionButton} onClick={clearAll}>AC</button>
                        <button className={styles.functionButton} onClick={handlePlusMinus}>+/-</button>
                        <button className={styles.functionButton} onClick={handlePercentage}>%</button>
                        <button
                            className={styles.operationButton}
                            onClick={() => performOperation('/')}
                        >÷
                        </button>
                    </div>

                    <div className={styles.buttonRow}>
                        <button className={styles.numberButton} onClick={() => inputDigit('7')}>7</button>
                        <button className={styles.numberButton} onClick={() => inputDigit('8')}>8</button>
                        <button className={styles.numberButton} onClick={() => inputDigit('9')}>9</button>
                        <button
                            className={styles.operationButton}
                            onClick={() => performOperation('*')}
                        >×
                        </button>
                    </div>

                    <div className={styles.buttonRow}>
                        <button className={styles.numberButton} onClick={() => inputDigit('4')}>4</button>
                        <button className={styles.numberButton} onClick={() => inputDigit('5')}>5</button>
                        <button className={styles.numberButton} onClick={() => inputDigit('6')}>6</button>
                        <button
                            className={styles.operationButton}
                            onClick={() => performOperation('-')}
                        >-
                        </button>
                    </div>

                    <div className={styles.buttonRow}>
                        <button className={styles.numberButton} onClick={() => inputDigit('1')}>1</button>
                        <button className={styles.numberButton} onClick={() => inputDigit('2')}>2</button>
                        <button className={styles.numberButton} onClick={() => inputDigit('3')}>3</button>
                        <button
                            className={styles.operationButton}
                            onClick={() => performOperation('+')}
                        >+
                        </button>
                    </div>

                    <div className={styles.buttonRow}>
                        <button
                            className={`${styles.numberButton} ${styles.zeroButton}`}
                            onClick={() => inputDigit('0')}
                        >0
                        </button>
                        <button className={styles.numberButton} onClick={inputDot}>.</button>
                        <button
                            className={styles.equalsButton}
                            onClick={calculateResult}
                        >=
                        </button>
                    </div>
                </div>

            </main>
            <footer className={styles.footer}>
                <div className={styles.copyright}>
                    <Image
                        aria-hidden
                        src="/copyright.svg"
                        alt="copyright"
                        width={16}
                        height={16}
                    />
                    2025
                </div>
                <a
                    href="https://github.com/Mredust"
                    target="_blank"
                    rel="noreferrer"
                >
                    <Image
                        aria-hidden
                        src="/github.svg"
                        alt="Github"
                        width={16}
                        height={16}
                    />
                    Mredust
                </a>
            </footer>
        </div>
    );
}
