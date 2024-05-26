import librosa
import numpy as np
import argparse
import os
from concurrent.futures import ThreadPoolExecutor

import soundfile as sf
import tempfile
import shutil

def convert_to_wav(input_file):
    # Создаем временный файл WAV
    temp_wav_file = tempfile.NamedTemporaryFile(suffix=".wav", delete=False)

    # Конвертируем в WAV
    try:
        data, samplerate = sf.read(input_file)
        sf.write(temp_wav_file.name, data, samplerate, format='WAV', subtype='PCM_16')
    except Exception as e:
        print("Ошибка при конвертации файла в WAV:", e)
        return None

    return temp_wav_file.name

def preprocess_and_encode(wav_file, n_mfcc=13, max_pad_len=400):
    # Конвертируем аудиофайл в WAV, если он не в этом формате
    if not wav_file.endswith(".wav"):
        wav_file = convert_to_wav(wav_file)
        if not wav_file:
            return None

    # Загрузка аудиофайла
    y, sr = librosa.load(wav_file, sr=None)

    # Извлечение признаков голоса (MFCC)
    mfcc = librosa.feature.mfcc(y=y, sr=sr, n_mfcc=n_mfcc)

    # Padding или trimming до max_pad_len
    if mfcc.shape[1] > max_pad_len:
        mfcc = mfcc[:, :max_pad_len]
    else:
        pad_width = max_pad_len - mfcc.shape[1]
        mfcc = np.pad(mfcc, pad_width=((0, 0), (0, pad_width)), mode='constant')

    return mfcc


def compare_files_threaded(wav_file_1, wav_file_2):
    with ThreadPoolExecutor(max_workers=2) as executor:
        future_1 = executor.submit(preprocess_and_encode, wav_file_1)
        future_2 = executor.submit(preprocess_and_encode, wav_file_2)
        mfcc_1 = future_1.result()
        mfcc_2 = future_2.result()

    if mfcc_1 is None or mfcc_2 is None:
        print("Ошибка в предобработке аудиофайлов.")
        return

    # Вычисление косинусного сходства между признаками голоса
    similarity = np.dot(mfcc_1.flatten(), mfcc_2.flatten()) / (np.linalg.norm(mfcc_1.flatten()) * np.linalg.norm(mfcc_2.flatten())) * 100

    print("{:.2f}".format(similarity))

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Сравнение двух аудиофайлов формата WAV и вывод процента сходства.")
    parser.add_argument("file1", type=str, help="Путь к первому аудиофайлу (формат WAV)")
    parser.add_argument("file2", type=str, help="Путь ко второму аудиофайлу (формат WAV)")
    args = parser.parse_args()

    wav_file_1 = args.file1
    wav_file_2 = args.file2

    compare_files_threaded(wav_file_1, wav_file_2)